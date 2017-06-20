import requests
from lxml import html
import time

i = 0

def issue_to_num(issue_text):
    digits = []
    for symbol in issue_text:
        if '1234567890'.find(symbol) != -1:
            digits.append(int(symbol))
    return(digits)

issues = 0
start = time.time()
with open("file.txt") as file:
    lines = sum(1 for line in open("file.txt"))
    for line in file.readlines():
        i += 1
        response = requests.get(line)
        parsed_body = html.fromstring(response.text)
        issue_text = ' '.join(parsed_body.xpath('//div[@class="contest-item-status contest-item-status-disputed"]/a/text()'))
        issue_num = 0 if issue_text == '' else sum(issue_to_num(issue_text))
        issues += issue_num

stop = time.time()
print("Total:", issues, "issues")
print("Execution time:", stop-start, "secs")