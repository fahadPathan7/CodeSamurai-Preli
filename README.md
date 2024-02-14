# CodeSamurai Preliminary Round 2024

## Question
[Preliminary round](./Questions/[Updated]%20Preliminary%20Round%20Problem%20Statement.pdf)

## Run the code
Use docker to run the code. The following command will build the docker image and run the code.

```bash
docker build --tag=sol:latest .
docker run -dit -p 8000:8000 --rm --name=sol sol:latest
```