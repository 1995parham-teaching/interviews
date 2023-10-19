class Node:
    def __init__(self, value, parent):
        self.value = value
        self.leaf = False
        self.children = {}
        self.parent = parent


def traverse(prefix: str, current_node: Node):
    if current_node.leaf:
        print(prefix)
    for child_value, child_node in current_node.children.items():
        # if child_node.leaf:
        #     print(prefix + child_value)
            # continue
        traverse(prefix + child_value, child_node)


def traverse_count(current_node: Node):
    sum = 0
    if current_node.leaf:
        sum += 1
    for child_value, child_node in current_node.children.items():
        # if child_node.leaf:
        #     sum += 1
            # continue
        sum += traverse_count(child_node)
    return sum


database_root = Node(value="", parent=None)

while True:
    command = input("Enter your command: ")

    if command == "insert":
        word = input("Enter your word: ")
        database_current_node = database_root
        for ch in word:
            if database_current_node.children.get(ch, 0) == 0:
                database_current_node.children[ch] = Node(ch, database_current_node)
            database_current_node = database_current_node.children.get(ch, 0)
        database_current_node.leaf = True

    else:
        prefix = input("Enter your prefix: ")
        database_current_node = database_root
        non_exist = False
        for ch in prefix:
            if database_current_node.children.get(ch, 0) == 0:
                print("No word with given prefix exist")
                non_exist = True
                break
            database_current_node = database_current_node.children.get(ch, 0)
        if non_exist:
            continue
        if command == "look up":
            traverse(prefix, database_current_node)
        elif command == "count":
            print(traverse_count(database_current_node))
        elif command == "delete":
            if database_current_node.leaf:
                database_current_node.parent.children.pop(database_current_node.value)
            else:
                database_current_node.children = {}
        else:
            print("This command doesn't exist")
