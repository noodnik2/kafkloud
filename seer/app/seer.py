import chromadb

from chromadb.config import Settings

from langchain.vectorstores import Chroma
from langchain.embeddings import OpenAIEmbeddings
from langchain.text_splitter import RecursiveCharacterTextSplitter
from langchain.llms import OpenAI
from langchain.chains import RetrievalQA
from langchain.schema import Document
from langchain.document_loaders import TextLoader


def seer_logger(message):
    print(message)


class Seer:

    def __init__(self, host="localhost", port="8000", log_fn=seer_logger):
        chroma_settings = Settings(chroma_api_impl="rest", chroma_server_host=host, chroma_server_http_port=port)
        self.chroma_client = chromadb.Client(chroma_settings)
        self.embeddings = OpenAIEmbeddings()
        self.log = log_fn

    def ask(self, questions):
        self.log(f"received questions({questions})")
        vectordb = Chroma(embedding_function=self.embeddings, client=self.chroma_client)
        qa = RetrievalQA.from_chain_type(
            llm=OpenAI(),
            chain_type="stuff",
            retriever=vectordb.as_retriever()
        )
        answers = []
        for q in questions:
            self.log(f"running query: {q}")
            answer = qa.run(q)
            self.log(f"the answer is: {answer}")
            answers.append(answer)
        return answers

    def load(self, fns):
        self.log(f"received fns({fns})")
        documents_to_split = []
        for fn in fns:
            self.log(f"reading({fn})")
            loader = TextLoader(fn)
            for doc in loader.load():
                documents_to_split.append(doc)
        self._split_and_load_documents(documents_to_split)

    def accept(self, texts):
        self.log(f"received texts({texts})")
        documents_to_split = []
        for text in texts:
            documents_to_split.append(Document(page_content=text))
        self._split_and_load_documents(documents_to_split)

    def _split_and_load_documents(self, documents_to_split):
        self.log("RecursiveCharacterTextSplitter")
        text_splitter = RecursiveCharacterTextSplitter(chunk_size=1000, chunk_overlap=0)
        self.log("text_splitter.split_documents")
        documents = text_splitter.split_documents(documents_to_split)
        self.log("Chroma.from_documents")
        Chroma.from_documents(documents, self.embeddings, client=self.chroma_client)

