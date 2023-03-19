export const config = {
  useMockData: false, // Set this to 'false' to use real data from the Golang backend
};

export const mockConversations = [
  {
    id: "conversation_0",
    name: "Poems",
    messages: [
      "Can you share a famous poem?",
      "Certainly! Here's a famous poem by Robert Frost called 'The Road Not Taken':\n\nTwo roads diverged in a yellow wood,\nAnd sorry I could not travel both\nAnd be one traveler, long I stood\nAnd looked down one as far as I could\nTo where it bent in the undergrowth;",
      "What's the meaning of the poem?",
      "The poem is often interpreted as a reflection on the choices we make in life and their impact on our personal journey. The speaker in the poem comes across two paths in the woods and chooses the one less traveled, symbolizing the decision to take a unique or unconventional route in life.",
      "Can you share another poem?",
      "Of course! Here's a short poem by Emily Dickinson called 'Hope is the thing with feathers':\n\nHope is the thing with feathers\nThat perches in the soul,\nAnd sings the tune without the words,\nAnd never stops at all.",
    ],
  },
  {
    id: "conversation_1",
    name: "Code Excerpts",
    messages: [
      "Can you give me a code example in Python?",
      "Sure, here's a simple Python function that calculates the factorial of a number using recursion:\n\n```\ndef factorial(n):\n    if n == 1:\n        return 1\n    else:\n        return n * factorial(n - 1)\n\nprint(factorial(5))\n```",
      "What does recursion mean?",
      "Recursion is a programming technique where a function calls itself as a subroutine, allowing it to break a complex problem into smaller, more manageable subproblems. It is often used as an alternative to iterative looping constructs, such as 'for' or 'while' loops.",
      "Can you give me an example of a loop in Python?",
      "Certainly! Here's a simple example of a 'for' loop in Python that prints the numbers from 1 to 5:\n\n```\nfor i in range(1, 6):\n    print(i)\n```",
      "How can I calculate the factorial using a loop?",
      "You can calculate the factorial using a loop like this:\n\n```\ndef factorial(n):\n    result = 1\n    for i in range(1, n + 1):\n        result *= i\n    return result\n\nprint(factorial(5))\n```",
    ],
  },
  {
    id: "conversation_2",
    name: "Food Recipes",
    messages: [
      "Can you suggest a simple pasta recipe?",
      "Of course! Here's a simple recipe for spaghetti aglio e olio:\n\nIngredients:\n- 400g spaghetti\n- 4 garlic cloves, thinly sliced\n- 1/2 cup olive oil\n- 1/2 teaspoon red pepper flakes\n- Salt\n- Freshly ground black pepper\n- Chopped parsley\n- Grated Parmesan\n\nInstructions:\n1. Cook the spaghetti in a large pot of salted boiling water until al dente.\n2. In a large pan, heat the olive oil over medium heat. Add the garlic and red pepper flakes, and cook until the garlic is golden brown.\n3. Drain the spaghetti and add it to the pan with the garlic and oil. Toss to coat the pasta.\n4. Season with salt and black pepper, and garnish with chopped parsley and grated Parmesan.",
      "What other pasta dishes do you recommend?",
      "I recommend trying pasta alla carbonara, which is a delicious pasta dish made with eggs, cheese, pancetta, and black pepper. Another great option is pasta primavera, a pasta dish made with fresh, sautéed vegetables and a light, creamy sauce.",
      "Can you suggest a vegetarian pasta recipe?",
      "Certainly! You can try a simple spinach and ricotta stuffed shells recipe:\n\nIngredients:\n- 12 jumbo pasta shells\n- 1 cup ricotta cheese\n- 1 cup chopped fresh spinach\n- 1/4 cup grated Parmesan\n- 1 egg\n- Salt and pepper\n- 2 cups marinara sauce\n\nInstructions:\n1. Cook the pasta shells according to the package instructions. Drain and set aside.\n2. In a bowl, mix the ricotta, spinach, Parmesan, egg, salt, and pepper.\n3. Preheat the oven to 350°F (180°C).\n4. Spread 1 cup of marinara sauce in the bottom of a baking dish.\n5. Stuff each pasta shell with the spinach and ricotta mixture, and place them in the baking dish.\n6. Pour the remaining marinara sauce over the stuffed shells.\n7. Bake for 25-30 minutes, or until the sauce is bubbling and the shells are heated through.",
    ],
  },
  {
    id: "conversation_3",
    name: "Movie Summaries",
    messages: [
      "Can you give me a summary of the movie 'Inception'?",
      "Inception (2010) is a science-fiction film directed by Christopher Nolan. The movie follows Dom Cobb (played by Leonardo DiCaprio), a skilled thief who steals information from people's minds while they are dreaming. Cobb is offered a chance to have his criminal history erased if he can successfully plant an idea in the mind of a corporate heir using a technique called inception. To accomplish this, Cobb assembles a team of specialists, and they delve into the complex world of dreams within dreams. Throughout the film, they face various obstacles as they try to complete their mission while battling their own inner demons.",
      "What is the significance of the spinning top?",
      "The spinning top is a totem used by Cobb to determine whether he is in a dream or in reality. If the top keeps spinning indefinitely, it indicates that he is still in a dream. If the top wobbles and falls, it means he is in reality. The ambiguity of the spinning top at the end of the movie leaves the audience wondering whether Cobb is still in a dream or has returned to reality.",
      "Can you recommend another movie by Christopher Nolan?",
      "I would recommend watching 'Interstellar' (2014), another science-fiction film directed by Christopher Nolan. The movie follows a group of astronauts who travel through a wormhole in search of a new habitable planet for humanity. The film explores themes of love, sacrifice, and the nature of time.",
    ],
  },
];

export const mockSystemodels = ["Model 1", "Model 2", "Model 3"];
