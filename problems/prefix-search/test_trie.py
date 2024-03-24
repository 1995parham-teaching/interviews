import unittest

from trie import Trie


class TestTrie(unittest.TestCase):
    def test_insert(self):
        t = Trie()
        t.insert("Parham")
        t.insert("Elahe")
        t.insert("Parya")
        t.insert("Elie")

        self.assertListEqual(t.traverse(""), ["Parham", "Parya", "Elahe", "Elie"])
        self.assertListEqual(t.traverse("Par"), ["Parham", "Parya"])
        self.assertListEqual(t.traverse("El"), ["Elahe", "Elie"])
        self.assertListEqual(t.traverse("Ela"), ["Elahe"])
        self.assertListEqual(t.traverse("Eli"), ["Elie"])
        self.assertListEqual(t.traverse("Ele"), [])

    def test_delete(self):
        t = Trie()
        t.insert("Parham")
        t.insert("Elahe")
        t.insert("Parya")
        t.insert("Elie")

        self.assertListEqual(t.traverse(""), ["Parham", "Parya", "Elahe", "Elie"])
        t.delete("Par")
        self.assertListEqual(t.traverse(""), ["Elahe", "Elie"])


if __name__ == "__main__":
    unittest.main()
