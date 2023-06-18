import sys
import argparse

from dotenv import load_dotenv
from seer import Seer

parser = argparse.ArgumentParser(
    prog=sys.argv[0],
    description="Loads and/or queries data into/from ChromaDB using OpenAI",
    allow_abbrev=True
)

parser.add_argument("-a", "--accept", help="accept text", type=str, action="append")
parser.add_argument("-l", "--load", help="load file", type=str, action="append")
parser.add_argument("-q", "--query", help="ask question", type=str, action="append")
parser.add_argument("-v", "--verbose", help="enable verbose mode", action="store_true")
parser.add_argument("-c", "--chromadb", help="host:port for chromadb", default="localhost:8000")
args = parser.parse_args()

have_files_to_load = args.load and len(args.load) > 0
have_texts_to_accept = args.accept and len(args.accept) > 0
have_queries = args.query and len(args.query) > 0

if not have_files_to_load and not have_texts_to_accept and not have_queries:
    print(f"no action(s) specified")
    parser.print_usage()
    exit(1)


def print_verbose(message):
    if args.verbose:
        print(message)


load_dotenv()

chromadb_host, chromadb_port = args.chromadb.split(":")
print_verbose(f"using chromadb({chromadb_host}:{chromadb_port})")

seer = Seer(host=chromadb_host, port=chromadb_port, log_fn=print_verbose)

if args.load:
    seer.load(args.load)

if args.accept:
    seer.accept(args.accept)

if args.query:
    seer.ask(args.query)
