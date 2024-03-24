class Node:
    def __init__(self):
        self.children: dict[str, Node] = {}
        self.word: str | None = None
        self.index: int | None = None

    def child(self, ch: str):
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
    def __insert(cls, root: Node, word: str, actual_word: str, actual_index):
        if len(word) == 0:
            root.word = actual_word
            root.index = actual_index
            return
        node = root.child(word[0])
        cls.__insert(node, word[1:], actual_word, actual_index)

    def insert(self, word: str, index: int):
        self.__insert(self.root, word, word, index)

    def traverse(self, prefix: str) -> list[Node]:
        return self.__traverse(self.root, prefix, [])

    @classmethod
    def __traverse(cls, root: Node, prefix: str, words: list[Node]) -> list[Node]:
        if len(prefix) == 0 and root.word is not None:
            words.append(root)
        if len(prefix) > 0:
            return cls.__traverse(root.child(prefix[0]), prefix[1:], words)
        for _, node in root.children.items():
            cls.__traverse(node, "", words)
        return words

    def delete(self, prefix: str):
        nodes = self.__traverse(self.root, prefix, [])
        for node in nodes:
            if node.word is not None:
                node.word = None


class WordFilter:
    def __init__(self, words: list[str]):
        self.pref_trie = Trie()
        self.suff_trie = Trie()

        for i, word in enumerate(words):
            self.pref_trie.insert(word, i)
            self.suff_trie.insert(word[::-1], i)

    def f(self, pref: str, suff: str) -> int:
        index = -1
        pref_nodes = self.pref_trie.traverse(pref)
        suff_nodes = self.suff_trie.traverse(suff[::-1])

        pref_indecies: set[int] = set()

        for node in pref_nodes:
            if node.index is not None:
                pref_indecies.add(node.index)

        for node in suff_nodes:
            if node.index in pref_indecies and node.index > index:
                index = node.index

        return index


if __name__ == "__main__":
    wf = WordFilter(["apple", "Parham"])
    assert wf.f("a", "e") == 0
    assert wf.f("Par", "ham") == 1
    assert wf.f("Par", "pam") == -1
