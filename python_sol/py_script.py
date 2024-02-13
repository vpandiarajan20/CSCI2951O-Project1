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
truth_values_str = "[1:false 2:false 3:false 4:false 5:false 6:true 7:false 8:false 9:false 10:false 11:false 12:false 13:false 14:false 15:false 16:false 17:true 18:false 19:false 20:false 21:false 22:false 23:true 24:false 25:false 26:false 27:false 28:false 29:false 30:false 31:false 32:false 33:false 34:false 35:false 36:false 37:false 38:false 39:true 40:false 41:false 42:false 43:false 44:true 45:false 46:false 47:false 48:false 49:false 50:false 51:false 52:true 53:false 54:false 55:false 56:false 57:false 58:false 59:false 60:false 61:false 62:false 63:false 64:false 65:false 66:false 67:false 68:false 69:false 70:true 71:true 72:false 73:false 74:false 75:false 76:false 77:false 78:false 79:false 80:false 81:false 82:false 83:false 84:false 85:false 86:false 87:false 88:true 89:false 90:false 91:false 92:false 93:false 94:false 95:true 96:false 97:false 98:false 99:false 100:false 101:true 102:false 103:false 104:false 105:false 106:false 107:false 108:false 109:false 110:false 111:false 112:false 113:true 114:false 115:false 116:false 117:false 118:false 119:false 120:false 121:false 122:false 123:false 124:true 125:false 126:false 127:false 128:false 129:false 130:false 131:false 132:true 133:false 134:false 135:false 136:false 137:false 138:true 139:false 140:false 141:false 142:false 143:false 144:false 145:false 146:false 147:false 148:false 149:false 150:false 151:false 152:false 153:true 154:false 155:false 156:false 157:false 158:false 159:false 160:false 161:false 162:false 163:true 164:false 165:false 166:false 167:false 168:false 169:false 170:true 171:false 172:false 173:false 174:false 175:true 176:false 177:false 178:false 179:false 180:false 181:false]"

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
    value = True if value == "true" else False
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
solve_dimacs('../src/SAT-Solver/input/C181_3151.cnf')