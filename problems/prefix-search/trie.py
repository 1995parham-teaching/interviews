from __future__ import annotations


class Node:
    def __init__(self):
        self.children: dict[str, Node] = {}
        self.word: str | None = None

    def child(self, ch: str) -> Node:
        if len(ch) != 1:
            raise ValueError("ch should have only one character length")

        node = self.children.get(ch, None)
        if node is None:
            node = Node()
            self.children[ch] = node
        return node


class Trie:
    def __init__(self):
        self.root = Node()

    @classmethod
    def __insert(cls, root: Node, word: str, actual_word: str):
        if len(word) == 0:
            root.word = actual_word
            return
        node = root.child(word[0])
        cls.__insert(node, word[1:], actual_word)

    def insert(self, word: str):
        self.__insert(self.root, word, word)

    def traverse(self, prefix: str) -> list[str]:
        return self.__traverse(self.root, prefix, [])

    @classmethod
    def __traverse(cls, root: Node, prefix: str, words: list[str]) -> list[str]:
        if len(prefix) == 0 and root.word is not None:
            words.append(root.word)
        if len(prefix) > 0:
            return cls.__traverse(root.child(prefix[0]), prefix[1:], words)
        for _, node in root.children.items():
            cls.__traverse(node, "", words)
        return words
