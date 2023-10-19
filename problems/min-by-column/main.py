import pandas as pd


def min_by_column(database, col_names):
    # Fill the null values with zero
    df = database.fillna(0)

    # Iterate over the list of column names by order
    for col_name in col_names:
        # Shrink the DF by the min value for col_name
        df = df[df[col_name] == df[col_name].min()]

        # I just one row has remained we can return it
        if len(df) == 1:
            return df.iloc[0].to_dict()

    # It iteration over all column names has been finished we can just return the first row
    return df.iloc[0].to_dict()


table_1 = pd.DataFrame([1, 2, 3], columns=["a"])
# print(min_by_column(table_1, ["a"]))
assert min_by_column(table_1, ["a"]) == {"a": 1}

table_2 = pd.DataFrame([[1, 2], [3, 0]], columns=["a", "b"])
# print(min_by_column(table_2, "b"))
assert min_by_column(table_2, "b") == {"a": 3, "b": 0}

table_3 = pd.DataFrame([[1, -2], [3]], columns=["a", "b"])
# print(min_by_column(table_3, "b"))
assert min_by_column(table_3, "b") == {"a": 1, "b": -2}


table_5 = pd.DataFrame([[1, 3], [1, 0]], columns=["x", "y"])
# print(min_by_column(table_5, ["x", "y"]))
assert min_by_column(table_5, ["x", "y"]) == {"x": 1, "y": 0}

table_6 = pd.DataFrame([[2, 3], [2, 1], [1, 10]], columns=["x", "y"])
# print(min_by_column(table_6, ["x", "y"]))
assert min_by_column(table_6, ["x", "y"]) == {"x": 1, "y": 10}


table_7 = pd.DataFrame([[3, -1, 0], [1, 10, 1], [1, 10, 0]], columns=["x", "y", "z"])
# print(min_by_column(table_7, ["x", "y", "z"]))
assert min_by_column(table_7, ["x", "y", "z"]) == {"x": 1, "y": 10, "z": 0}


table_8 = pd.DataFrame([[1, 2, 3], [2, 2, 2]], columns=["x", "y", "z"])
# print(min_by_column(table_8, ["x", "y", "z"]))
assert min_by_column(table_8, ["x", "y", "z"]) == {"x": 1, "y": 2, "z": 3}

