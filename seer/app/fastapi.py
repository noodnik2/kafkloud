
from os import getenv
import logging
import asyncio

from fastapi import FastAPI, HTTPException
from starlette.requests import Request
from contextlib import asynccontextmanager
from threading import Thread

from app.seer import Seer
from app.kafka import KafkaSeer


class AIOKafkaSeer:
    def __init__(self, seer, broker):
        self._log_handler = logging.getLogger(__name__)
        self._seer = seer
        self._broker = broker
        self._loop = asyncio.get_event_loop()
        self._log_handler.debug(f"using kafka_broker({self._broker})")
        self._kafka_seer = KafkaSeer(self._seer, self._broker)
        self._poll_thread = Thread(target=self._poll_loop)
        self._poll_thread.start()
        self._log_handler.debug("launched")

    def close(self):
        self._log_handler.debug("closing")
        self._kafka_seer.cancel()
        self._poll_thread.join()
        self._log_handler.debug("thread joined")

    def _poll_loop(self):
        self._log_handler.debug(f"starting polling")
        self._kafka_seer.run()
        self._log_handler.debug(f"ended polling")


@asynccontextmanager
async def lifespan(app: FastAPI):
    app.state.logger = logging.getLogger(__name__)
    app.state.logger.info("setting DEBUG logging level")
    logging.basicConfig(format='%(asctime)s %(levelname)s: %(message)s', level=logging.DEBUG)
    chroma_db = getenv("CHROMA_DB", "localhost:8000")
    kafka_broker = getenv("KAFKA_BROKER")
    if not chroma_db or not kafka_broker or chroma_db == "" or kafka_broker == "":
        raise "incomplete configuration - missing 'CHROMA_DB' and/or 'KAFKA_BROKER'"
    chromadb_host, chromadb_port = chroma_db.split(":")
    app.state.logger.debug(f"using chromadb({chromadb_host}:{chromadb_port})")
    app.state.seer = Seer(host=chromadb_host, port=chromadb_port)
    aio_kafka_seer = AIOKafkaSeer(app.state.seer, kafka_broker)
    app.state.logger.debug("app initialized")
    yield
    app.state.logger.debug("app shutdown started")
    aio_kafka_seer.close()
    app.state.logger.debug("app shutdown completed")


app = FastAPI(lifespan=lifespan)


def bytes_to_lines(body_bytes) -> list:
    body_str = body_bytes.decode("utf-8")
    return body_str.split("\n")


@app.get("/_health")
async def health():
    return {"health": "OK"}


@app.post("/load")
async def load(r: Request):
    filenames = bytes_to_lines(await r.body())
    app.state.seer.load(filenames)
    return {"status": f"{len(filenames)} file(s) loaded"}


@app.post("/accept")
async def accept(r: Request):
    statement = (await r.body()).decode("utf-8")
    app.state.seer.accept([statement])
    return {"status": f"statement accepted"}


@app.get("/ask")
async def ask(r: Request):
    question = (await r.body()).decode("utf-8")
    answers = app.state.seer.ask([question])
    if len(answers) != 1:
        return {"error": f"{len(answers)} answer(s) unexpectedly returned"}
    return {"question": question, "answer": answers[0]}
