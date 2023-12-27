from look_up import Database


database = Database()

while True:
    command = input("Enter your command: ")

    if command == "insert":
        word = input("Enter your word: ")
        database.insert(word)
    else:
        prefix = input("Enter your prefix: ")
        current_node = database.get_prefix(prefix)
        if current_node is None:
            continue
        if command == "look up":
            database.traverse(prefix, current_node)
        elif command == "count":
            print(database.traverse_count(current_node))
        elif command == "delete":
            database.delete(current_node)
        else:
            print("This command doesn't exist")
