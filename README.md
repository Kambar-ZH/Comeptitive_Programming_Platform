# Competitive programming site


# client_side
1️⃣ Each folder in this directory indicates problem id and contains solutions for the corresponding problem.

# cmd
2️⃣ Contains: 
 - main.go - entry point of the program
 - 2 go / exe files (participant_solution.go / participant_solution.exe and main_solution.go / main_solution.exe)
 - makefile that builds binaries for the participant_solution.go and main_solution.go

# internal
3️⃣ Server and the logic of handling requests.

# temp_solutions
4️⃣ Solutions from participants are temporarly stored here.

# test
5️⃣ Encapsulates the logic of compiling participant_solution.go and main_solution.go. 
 - write new solutions into corresponding go files
 - execute makefile, which generates binaries for the solutions
 - compare the outputs of 2 programs
 - contains authors solutions with the test files.

# web
6️⃣ The frontend part (templates and css) are located here. 
