# 13. Design Improvements to Your Machine to Get Around Problems

## Core Principle
```
Build the machine, don't just do the work
Systemize solutions → Prevent recurrence
```

## Design Before People

### Sequence
```
1. What outcomes needed?
2. What machine produces those?
3. What roles required?
4. What people for those roles?
```

### Don't Design Around People
- If role needs X, role needs X
- Don't reshape role for person's weakness
- Find right person for designed role

## Movie Script Visualization

### Technique
```
Visualize the solution running like a movie
- Who does what?
- In what sequence?
- What are the handoffs?
- Where could it fail?
```

### Test the Script
- Walk through step by step
- Identify dependencies
- Find potential failures
- Build in checkpoints

## Iterative Process

### Design Loop
```
Diagnose → Design → Implement → Observe → Adjust → Repeat
```

### Quality Day Innovation
- Periodically stop and review
- Challenge existing design
- Look for improvements
- Implement best ideas

## Design Principles

### 1. Clear Accountability
```
Every task → One responsible person
No ambiguity about who owns what
```

### 2. Pyramid Structure
- Clear reporting lines
- Authority matches responsibility
- Escalation paths defined

### 3. Checks and Balances
- No single point of failure
- Verification built in
- Guardrails for high-stakes

### 4. Right Granularity
| Too Coarse | Too Fine |
|------------|----------|
| Missing control | Bureaucracy |
| Things fall through | Slow execution |
| Accountability unclear | Over-management |

## Designing Around Weaknesses

### People Weaknesses
```
Option 1: Train (if skill gap)
Option 2: Guardrail (add checks)
Option 3: Move (different role)
Option 4: Remove (if fundamental)
```

### System Weaknesses
```
Option 1: Add process step
Option 2: Add automation
Option 3: Add verification
Option 4: Redesign entirely
```

## Common Design Patterns

### Redundancy
- Backup for critical functions
- Cross-training
- Documentation

### Feedback Loops
- Metrics that trigger action
- Regular review cycles
- Automatic alerts

### Escalation
- Clear triggers
- Defined paths
- Time limits

## Avoid Over-Design

### Signs of Over-Design
- Too many steps
- Too much documentation
- Slows everything down
- People work around it

### Balance
```
Enough control to prevent problems
Not so much it prevents progress
```

## Key Formulas
| Formula | Meaning |
|---------|---------|
| Design BEFORE People | Sequence matters |
| One task = One owner | Clear accountability |
| Diagnose → Design → Implement | Iterative |
| Control ≈ Speed | Balance needed |
