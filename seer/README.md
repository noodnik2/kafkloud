# `seer` - _one that sees_

Seer can tell you things you may not have realized about what's
inherently being communicated in a set of data.

## Main Idea

Get on-board the AI craze and insert a GPT-like "seer" into the Kafkloud
data pipeline.

The main use case is to "teach" Seer things that it can later recall,
summarize and put into a greater "context" in response to questions
about what it's learned.

## Roadmap

### Current State
As of this writing, `seer` implements a [Python module](app/seer.py)
for which there's a [CLI](app/cli.py) and [HTTP](app/fastapi.py) interface,
used now mainly for testing.  Its packaging includes standalone and Docker
deployments, and its Kubernetes deployment artifacts are partially implemented.
See the [Makefile](Makefile) for the related deployment targets.

### Future State
It's envisioned to add interface(s) to `seer` enabling it to respond to
incoming questions on a (Kafka) topic, and deliver answers to another.


## Dialogues

### Using the CLI
```shell
$ python app/cli.py -v -q "Do you know anything about the Wilymajinkas?"
using chromadb(localhost:8000)
received questions(['Do you know anything about the Wilymajinkas?'])
running query: Do you know anything about the Wilymajinkas?
the answer is:  No, I don not know anything about the Wilymajinkas.
$ python app/cli.py -v -a "The Wilymajinkas are a tribe of northeastern Native Americans."  -a "They are very tame and should not be feared."
using chromadb(localhost:8000)
received texts(['The Wilymajinkas are a tribe of northeastern Native Americans.', 'They are very tame and should not be feared.'])
RecursiveCharacterTextSplitter
text_splitter.split_documents
Chroma.from_documents
$ python app/cli.py -v -q "Do you know anything about the Wilymajinkas?"
using chromadb(localhost:8000)
received questions(['Do you know anything about the Wilymajinkas?'])
running query: Do you know anything about the Wilymajinkas?
the answer is:  Yes, the Wilymajinkas are a tribe of northeastern Native Americans. They are very tame and should not be feared.
```

### Using HTTP

- [Seer Tests](./tests/seer_tests.http)