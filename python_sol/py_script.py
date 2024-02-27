from pysat.solvers import Solver
from pysat.formula import CNF

def solve_dimacs(dimacs_file):
    # Load DIMACS file
    with open(dimacs_file, 'r') as f:
        cnf_formula = CNF(from_fp=f)
    
    # Initialize solver
    solver = Solver(name='g4')
    
    # Add clauses to the solver
    for clause in cnf_formula.clauses:
        solver.add_clause(clause)
    
    # Solve the SAT problem
    sat = solver.solve()
    
    if sat:
        print("SATISFIABLE")
        # Print the solution
        solution = solver.get_model()
        print("Solution:", solution)
    else:
        print("UNSATISFIABLE")

from pysat.solvers import Solver

def verify_satisfiability(dimacs_file, truth_assignment):
    with open(dimacs_file, 'r') as f:
        cnf_formula = CNF(from_fp=f)
    # Initialize solver
    solver = Solver(name='g4')

    # Add clauses to the solver
    for clause in cnf_formula:
        solver.add_clause(clause)

    # Add assumptions based on the truth assignment
    for var, value in truth_assignment.items():
        if value:
            solver.add_clause([var])
        else:
            solver.add_clause([-var])
    
    # Check if the formula is satisfiable
    return solver.solve()

# Given truth assignment
# truth_assignment = {
#     1: False, 2: True, 3: True, 4: False, 5: True, 6: True, 7: True, 8: False,
#     9: False, 10: True, 11: False, 12: True, 13: True, 14: False, 15: True,
#     16: True, 17: True, 18: True, 19: False, 20: False, 21: True, 22: False,
#     23: False, 24: False, 25: False, 26: False
# }


# truth_assignment = {
#     1: True, 2: False, 3: True, 4: False, 5: False}

# Given string
truth_values_str = "final truth values map[0:2 1:2 2:2 3:2 4:2 5:2 6:2 7:2 8:2 9:1 10:2 11:2 12:1 13:2 14:2 15:2 16:2 17:2 18:2 19:2 20:2 21:2 22:2 23:2 24:2 25:2 26:2 27:2 28:2 29:2 30:1 31:1 32:2 33:2 34:2 35:2 36:2 37:2 38:2 39:2 40:2 41:2 42:2 43:2 44:2 45:2 46:2 47:1 48:2 49:2 50:2 51:2 52:2 53:2 54:1 55:2 56:2 57:2 58:2 59:2 60:2 61:2 62:2 63:2 64:2 65:2 66:2 67:2 68:1 69:2 70:2 71:2 72:2 73:1 74:2 75:2 76:2 77:2 78:2 79:2 80:2 81:2 82:2 83:2 84:2 85:1 86:2 87:2 88:2 89:2 90:2 91:2 92:2 93:2 94:2 95:2 96:1 97:2 98:2 99:2 100:2 101:2 102:2 103:2 104:2 105:2 106:2 107:1 108:2 109:2 110:2 111:2 112:2 113:2 114:2 115:2 116:2 117:1 118:2 119:2 120:2 121:2 122:2 123:2 124:2 125:2 126:2 127:1 128:2 129:2 130:2 131:2 132:2 133:1 134:2 135:2 136:2 137:2 138:2 139:1 140:2 141:2 142:2 143:2 144:2 145:2 146:2 147:2 148:2 149:1 150:2 151:2 152:2 153:2 154:2 155:2 156:2 157:2 158:2 159:1 160:2 161:2 162:2 163:2 164:2 165:1 166:2 167:2 168:2 169:2 170:2 171:2 172:2 173:1 174:2 175:2 176:2 177:2 178:2 179:2 180:2 181:2]"
# with open("output_assignments.txt", "r") as f:
#     truth_values_str = f.read().strip()
# Extracting the substring containing key-value pairs
start_index = truth_values_str.find("[")
end_index = truth_values_str.rfind("]")
key_value_pairs_str = truth_values_str[start_index + 1:end_index]

# Splitting the key-value pairs
key_value_pairs = key_value_pairs_str.split()

# Creating the dictionary
truth_assignment = {}
for pair in key_value_pairs:
    key, value = pair.split(":")
    key = int(key)
    if key == 0:
        continue
    value = True if value == "1" else False
    truth_assignment[key] = value


# if verify_satisfiability('../src/SAT-Solver/toy_lecture.cnf', truth_assignment):
# if verify_satisfiability('../src/SAT-Solver/toy_solveable.cnf', truth_assignment):
if verify_satisfiability('../src/SAT-Solver/input/C181_3151.cnf', truth_assignment):
    print("The formula is satisfiable with the given truth assignment.")
else:
    print("The formula is unsatisfiable with the given truth assignment.")

# # Example usage
# solve_dimacs('../src/SAT-Solver/toy_solveable.cnf')
# solve_dimacs('../src/SAT-Solver/toy_lecture.cnf')
# solve_dimacs('../src/SAT-Solver/input/C181_3151.cnf')
# solve_dimacs('../src/SAT-Solver/input/C208_3254.cnf')
# solve_dimacs('../src/SAT-Solver/input/C1065_082.cnf')
# solve_dimacs('../src/SAT-Solver/toy_lecture.cnf')