FROM python:3.9

WORKDIR /code

COPY ./requirements.txt /code/requirements.txt

RUN pip install -r /code/requirements.txt

COPY ./app /code/app

# TODO: remove assets from image; it's being used to support test use cases
COPY ./assets /code/assets

CMD ["uvicorn", "app.fastapi:app", "--host", "0.0.0.0", "--port", "8030"]
