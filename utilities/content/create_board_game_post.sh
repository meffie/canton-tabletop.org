#!/bin/bash 

# Script defaults 
min_article_words=300 
max_article_words=600 
create_as_draft="false" #"true" or "false" (printed in post metadata)

# Provide usage instructions if no input was provided. 
if [ $# -eq 0 ]; then
    echo "Usage: $0 <board-game>"
    exit 1
fi

# Description of the article provided by the user. 
board_game=$1

# Generate article metadata 
generated_title=$(ollama run codellama:7b "Write a title \ 
    for a blog article about the board game '$board_game'. The article will be \
    displayed with other blog posts in the category: board games")
title=$(echo "$generated_title" | tr -d '\n' | sed -e 's/[[:space:]]*$//') # remove newlines and trim spaces. 

generated_description=$(ollama run codellama:7b "Write a description \ 
    for a blog article about the board game '$board_game' titled '$title'. The article will be \
    displayed with other blog posts in the category: board games")
description=$(echo "$generated_description" | tr -d '\n' | sed -e 's/[[:space:]]*$//') 

article_date=$(date "+%Y-%m-%dT%H:%M:%S%z") 

slug=$(echo "$title" | tr '[:upper:]' '[:lower:]' | sed -e 's/[^a-z0-9]/-/g' -e 's/^-//' -e 's/-$//')

echo "---"
echo "title: $title"
echo "date: $article_date"
echo "slug: /$slug"
echo "description: $description"
echo "image: images/catan-close-up.png"
echo "categories: "
echo "  - board games"
echo "draft: $create_as_draft"
echo "---"

# Generate the article contents. 
# ollama run codellama:7b "Write a blog article using markdown language (.md) \
#     The article should be between $min_article_words and $max_article_words words long. \
#     The largest markdown header should be size 2. \ 
#     The article is about $title \ 
#     A description of the article is $description "
