class Node:
    def __init__(self, value, parent):
        self.value = value
        self.leaf = False
        self.children = {}
        self.parent = parent


class Database:
    def __init__(self):
        self.root = Node(value="", parent=None)

    def insert(self, word):
        database_current_node = self.root
        for ch in word:
            if database_current_node.children.get(ch, 0) == 0:
                database_current_node.children[ch] = Node(value=ch, parent=database_current_node)
            database_current_node = database_current_node.children.get(ch, 0)
        database_current_node.leaf = True

    def traverse(self, prefix: str, current_node: Node):
        if current_node.leaf:
            print(prefix)
        for child_value, child_node in current_node.children.items():
            # if child_node.leaf:
            #     print(prefix + child_value)
            # continue
            self.traverse(prefix + child_value, child_node)

    def traverse_count(self, current_node: Node):
        sum = 0
        if current_node.leaf:
            sum += 1
        for child_value, child_node in current_node.children.items():
            # if child_node.leaf:
            #     sum += 1
            # continue
            sum += self.traverse_count(child_node)
        return sum

    def get_prefix(self, prefix: str):
        database_current_node = self.root
        for ch in prefix:
            if database_current_node.children.get(ch, 0) == 0:
                print("No word with given prefix exist")
                return None
            database_current_node = database_current_node.children.get(ch, 0)

        return database_current_node

    def delete(self, current_node: Node):
        if current_node.leaf:
            current_node.parent.children.pop(current_node.value)
        else:
            current_node.children = {}
