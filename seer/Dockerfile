FROM python:3.10-slim-bullseye

COPY ./requirements.txt /code/requirements.txt

RUN pip install --no-cache-dir -r /code/requirements.txt

COPY ./app /code/app

WORKDIR /code
CMD ["uvicorn", "app.fastapi:app", "--host", "0.0.0.0", "--port", "8030"]
