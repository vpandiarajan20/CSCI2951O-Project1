		go back until unique implication points - have one more literal from the current decision level
		1. have a conflict clause - that is initial unit clause
		2. look at the decision stack - pop off elements - even the latest decision that caused the conflict
		3. if not related to clause ignore it (unassign it); if decision (branch/implication) involves in clause,
		fetch all the literals that are parents of it --> of these, all the ones from previous decision levels go into clause
		all from previous decision level go into clasuse
		ones from current decision level added into queue -- only on the current queue left (might still have many from current decision level)
		when conflict then invert everything
		4. when current set has only 1 literal, then you stop - then you have conflict clause
		all vars from prev decision levels

		()
		while queue len > 1 {
			pop from deicion stack
			unassign unrelated to conflict
			related to conflict then fetch all parents and unassign -- just need immediate parents
			all parents from previous decision levels go into conflict cause, all parents from the current decision level go into queue
		}
        backtrack up to the max level of everything in previous decision levels
		
    once do this, generate a new unit want to do this properly
		backtrack all the way to level 0 if learn a unit clause

		1. first check if all stuff correct -> if step through this with debugger is useless -> have assertions that check that everything correct is dontconflictClause
		-> after conflict clause, everything except one variable is assigned and in conflict
		-> watched literals - after make changes then watched literals
		-> assertions catch bug
			-> if have time

		when backtrack, find the highest decision in the clause that is not from the current decision level
		then backtrack to that level
		the literal from current decision level becomes a unit - all other things are still assigned from clause
		every time create conflict clause, get a unit propagation immediately after