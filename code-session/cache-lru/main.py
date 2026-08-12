"""LRU cache with O(1) get and put.

A hash map gives O(1) access to an entry, and a doubly linked list keeps the
entries ordered from most recently used to least recently used so that eviction
is O(1) as well. Each node carries its own key: without it, evicting the node at
the tail would leave no way to find the map entry that points at it.
"""

from collections.abc import Iterator


class Node[K, V]:
    """An entry of the cache, and at the same time a node of the recency list."""

    __slots__ = ("key", "next", "prev", "value")

    def __init__(self, key: K, value: V) -> None:
        self.key: K = key
        self.value: V = value
        self.prev: Node[K, V] | None = None
        self.next: Node[K, V] | None = None


class RecencyList[K, V]:
    """A doubly linked list ordered from most recently used to least recently used.

    Every operation is O(1) and none of them searches the list: the caller always
    hands over the node it wants to move or remove.
    """

    def __init__(self) -> None:
        self.head: Node[K, V] | None = None  # most recently used
        self.tail: Node[K, V] | None = None  # least recently used

    def push_front(self, node: Node[K, V]) -> None:
        node.prev = None
        node.next = self.head

        if self.head is not None:
            self.head.prev = node
        self.head = node

        if self.tail is None:
            self.tail = node

    def remove(self, node: Node[K, V]) -> None:
        if node.prev is not None:
            node.prev.next = node.next
        else:
            self.head = node.next

        if node.next is not None:
            node.next.prev = node.prev
        else:
            self.tail = node.prev

        node.prev = None
        node.next = None

    def move_to_front(self, node: Node[K, V]) -> None:
        if self.head is node:
            return

        self.remove(node)
        self.push_front(node)

    def pop_back(self) -> Node[K, V] | None:
        node = self.tail
        if node is not None:
            self.remove(node)

        return node

    def __iter__(self) -> Iterator[Node[K, V]]:
        node = self.head

        while node is not None:
            yield node
            node = node.next


class LRUCache[K, V]:
    """A cache of at most `capacity` entries that evicts the least recently used one.

    Both `get` and `put` count as a use, including a `put` that overwrites an
    existing key.
    """

    def __init__(self, capacity: int) -> None:
        if capacity < 1:
            raise ValueError("capacity must be at least 1")

        self.capacity = capacity

        self._entries: dict[K, Node[K, V]] = {}
        self._recency: RecencyList[K, V] = RecencyList()

    def __len__(self) -> int:
        return len(self._entries)

    def __contains__(self, key: object) -> bool:
        """Membership does not count as a use, so it never reorders the cache."""
        return key in self._entries

    def __getitem__(self, key: K) -> V:
        """Return the value of `key`, or raise `KeyError` when it is not cached."""
        node = self._entries.get(key)
        if node is None:
            raise KeyError(key)

        self._recency.move_to_front(node)

        return node.value

    def __setitem__(self, key: K, value: V) -> None:
        self.put(key, value)

    def get[D](self, key: K, default: D) -> V | D:
        """Return the value of `key`, or `default` when it is not cached.

        There is no in-band value for a miss on purpose: `None` is a perfectly
        good thing to store, so a cache that answers `None` cannot tell "absent"
        from "present and empty". Callers that have no sensible default should
        use `cache[key]` and handle `KeyError`.
        """
        try:
            return self[key]
        except KeyError:
            return default

    def put(self, key: K, value: V) -> None:
        node = self._entries.get(key)

        # An overwrite is a use, not an insertion: the existing node moves to the
        # front and nothing is evicted, because the entry count does not grow.
        if node is not None:
            node.value = value
            self._recency.move_to_front(node)

            return

        if len(self._entries) >= self.capacity:
            self._evict()

        node = Node(key, value)
        self._entries[key] = node
        self._recency.push_front(node)

    def keys(self) -> list[K]:
        """The cached keys, most recently used first."""
        return [node.key for node in self._recency]

    def _evict(self) -> None:
        node = self._recency.pop_back()
        if node is None:
            return

        del self._entries[node.key]


if __name__ == "__main__":
    cache = LRUCache[int, int](2)

    cache.put(1, 1)
    cache.put(2, 2)
    assert cache.get(1, -1) == 1

    cache.put(3, 3)  # evicts key 2, the least recently used one
    assert cache.get(2, -1) == -1

    cache.put(4, 4)  # evicts key 1
    assert cache.get(1, -1) == -1
    assert cache.get(3, -1) == 3
    assert cache.get(4, -1) == 4

    # Overwriting a key is a use: it moves the key to the front, keeps a single
    # node for it, and evicts nobody.
    cache = LRUCache[int, int](2)
    cache.put(1, 1)
    cache.put(2, 2)
    cache.put(1, 10)
    assert cache.keys() == [1, 2]
    assert len(cache) == 2

    cache.put(3, 3)  # evicts key 2, because key 1 was just written
    assert cache.keys() == [3, 1]
    assert cache.get(1, -1) == 10
    assert cache.get(2, -1) == -1

    # A miss is not confused with a stored `None`.
    maybe = LRUCache[str, None](1)
    maybe.put("a", None)
    assert maybe.get("a", "missing") is None
    assert maybe.get("b", "missing") == "missing"
    assert "a" in maybe
    try:
        maybe["b"]
    except KeyError:
        pass
    else:
        raise AssertionError("a miss must raise KeyError")

    # Capacity 1: every insertion evicts the previous entry.
    tiny = LRUCache[str, str](1)
    tiny.put("a", "1")
    tiny.put("a", "2")
    assert tiny.keys() == ["a"]
    tiny.put("b", "1")
    assert tiny.keys() == ["b"]

    try:
        LRUCache[str, str](0)
    except ValueError:
        pass
    else:
        raise AssertionError("capacity must be validated")
