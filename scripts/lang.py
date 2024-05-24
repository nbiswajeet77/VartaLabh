#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
import openai
from langchain_openai import OpenAI
from langchain.prompts import PromptTemplate
from langchain.chains import LLMChain

# Set up your OpenAI API key
openai.api_key = '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'

class Node:
    def __init__(self, name, prompt_template, condition):
        self.name = name
        self.prompt_template = prompt_template
        self.condition = condition
        self.message_count = 0
        self.messages = []
        self.next_node = None

    def set_next(self, next_node):
        self.next_node = next_node

    def run_chain(self, conversation):
        llm = OpenAI(api_key=openai.api_key, temperature=0.7)
        chain = LLMChain(llm=llm, prompt=self.prompt_template)
        return chain.run(conversation=conversation)

# Define the prompt templates for each node
identify_mood_prompt = PromptTemplate(
    input_variables=["conversation"],
    template="Identify the user's mood based on the following conversation:\n{conversation}\nRespond in a supportive and understanding manner."
)

identify_issue_prompt = PromptTemplate(
    input_variables=["conversation"],
    template="Help the user identify their main issue based on the following conversation:\n{conversation}\nRespond in a thoughtful and guiding manner."
)

understand_negative_patterns_prompt = PromptTemplate(
    input_variables=["conversation"],
    template="Explore the user's negative thought patterns based on the following conversation:\n{conversation}\nRespond in an insightful and empathetic manner."
)

reshape_thought_patterns_prompt = PromptTemplate(
    input_variables=["conversation"],
    template="Assist the user in reshaping their thought patterns based on the following conversation:\n{conversation}\nRespond in an encouraging and constructive manner."
)

# Define the conditions for moving to the next node
def identify_mood_condition(node):
    return node.message_count >= 5

def identify_issue_condition(node):
    return node.message_count >= 6 or any("makes sense" in msg.lower() for msg in node.messages)

def understand_negative_patterns_condition(node):
    return node.message_count >= 7 or any("let's move ahead" in msg.lower() for msg in node.messages)

def reshape_thought_patterns_condition(node):
    return node.message_count >= 9

# Define the nodes with specific prompts
identify_mood_node = Node(
    name="Identify Mood",
    prompt_template=identify_mood_prompt,
    condition=identify_mood_condition
)
identify_issue_node = Node(
    name="Identify Issue",
    prompt_template=identify_issue_prompt,
    condition=identify_issue_condition
)
understand_negative_patterns_node = Node(
    name="Understand Negative Thought Patterns",
    prompt_template=understand_negative_patterns_prompt,
    condition=understand_negative_patterns_condition
)
reshape_thought_patterns_node = Node(
    name="Reshape Thought Patterns",
    prompt_template=reshape_thought_patterns_prompt,
    condition=reshape_thought_patterns_condition
)

# Define the transitions
identify_mood_node.set_next(identify_issue_node)
identify_issue_node.set_next(understand_negative_patterns_node)
understand_negative_patterns_node.set_next(reshape_thought_patterns_node)

class Router:
    def __init__(self, start_node):
        self.current_node = start_node

    def handle_message(self, user_input):
        self.current_node.messages.append(user_input)
        self.current_node.message_count += 1

        # Get the response from the LLM using the current node's specific prompt
        bot_response = self.current_node.run_chain(self.current_node.messages)
        self.current_node.messages.append(bot_response)
        self.current_node.message_count += 1
        print(bot_response)

        # Check if the current node's condition is met to move to the next node
        if self.current_node.condition(self.current_node):
            self.current_node = self.current_node.next_node
            if self.current_node:
                print(f"Transitioning to {self.current_node.name}")
            else:
                print("Conversation flow completed.")
                return

# Initialize the router with the start node
router = Router(start_node=identify_mood_node)

while True:
    user_input = input("Enter a message (or type 'exit' to stop): ")
    if user_input.lower() == 'exit':
        break
    router.handle_message(user_input)