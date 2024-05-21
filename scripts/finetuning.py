from openai import OpenAI
client = OpenAI(api_key="enter-key-here")

client.fine_tuning.jobs.create(
  training_file="file-zTfK21gpZW1ZfVEMKJmT30Pb", 
  model="gpt-3.5-turbo"
)