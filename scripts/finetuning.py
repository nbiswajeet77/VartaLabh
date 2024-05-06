from openai import OpenAI
client = OpenAI(api_key="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")

client.fine_tuning.jobs.create(
  training_file="file-zTfK21gpZW1ZfVEMKJmT30Pb", 
  model="gpt-3.5-turbo"
)