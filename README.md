# Exercise Generator

## Run API server in local

```
go run cmd/*.go server
```

## Run evaluation script

### Step 1: Specify input/output dir with your own config.yaml

- Copy and fill value in your own `config.yaml` file
- Specify the text file that contains all your words to evaluate the model with field `EVALUATION_CONFIG.BASELINE_FILE_PATH`
- Specify the result directory with field `EVALUATION_CONFIG.RESULT_DIR`

### Step 2: Run command

```
go run cmd/*.go generate-evaluation-template -questionType=<question_type> -provider=openai
```

For `<question_type>`, there are 3 options:

1. `definition`: to generate and evaluate the language model with **definition** question
2. `synonym`: to generate and evaluate the language model with **synonym** question
3. `application`: to generate and evaluate the language model with **application** question
