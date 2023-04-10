with open('cleaned_data.txt', 'w') as out:
    with open('unformatted_data.txt', 'r') as f:
        for i, line in enumerate(f):
            if i == 0:
                num_cols = len(line.split(","))
                out.write(line)
                continue

            current_line = line.strip("\n").split(",")
            print(current_line)

            current_line = current_line[:2] + [x.strip("dab*") for i, x in enumerate(current_line) if i>1]
            
            print(current_line)
            new_line = current_line + ([""] * (num_cols - len(current_line)))
            print(num_cols - len(current_line))
            out.write(",".join(new_line)+"\n")
