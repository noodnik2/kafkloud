from os import getenv
from fastapi import FastAPI
from starlette.requests import Request
from app.seer import Seer


def log(message):
    print(message)


async def request_list(r: Request) -> list:
    body_bytes = await r.body()
    body_str = body_bytes.decode("utf-8")
    return body_str.split("\n")


chromadb = getenv("CHROMADB", "localhost:8000")
chromadb_host, chromadb_port = chromadb.split(":")
log(f"using chromadb({chromadb_host}:{chromadb_port})")

app = FastAPI()
seer = Seer(host=chromadb_host, port=chromadb_port, log_fn=log)


@app.get("/_health")
async def health():
    return {"health": "OK"}


@app.get("/ask")
async def ask(r: Request):
    questions = await request_list(r)
    answers = seer.ask(questions)
    return {"answers": answers}


@app.post("/load")
async def load(r: Request):
    filenames = await request_list(r)
    seer.load(filenames)
    return {"status": f"{len(filenames)} file(s) loaded"}


@app.post("/accept")
async def accept(r: Request):
    texts = await request_list(r)
    seer.accept(texts)
    return {"status": f"{len(texts)} text(s) accepted"}
