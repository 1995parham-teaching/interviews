import math

import numpy as np
import pandas as pd


def dist(X: pd.Series, Y: pd.Series):
    return math.dist(X, Y)


class KNN:
    def __init__(self, database: pd.DataFrame, k: int):
        self.database = database
        self.k = k

    def predict(self, req: pd.Series) -> pd.DataFrame:
        self.database["distance"] = self.database.apply(lambda x: dist(x, req), axis=1)
        k_nearest_neighbours = self.database.iloc[
            np.argpartition(self.database["distance"], self.k - 1)[: self.k]
        ]

        return k_nearest_neighbours


if __name__ == "__main__":
    data = pd.DataFrame([[1, 2], [3, 2], [5, 4], [3, 1], [6, 3], [8, 4]])

    knn = KNN(data, 3)

    req = pd.Series([5, 5])
    print(knn.predict(req))
