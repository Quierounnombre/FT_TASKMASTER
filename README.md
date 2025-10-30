## How to develop

Start up the environment:
'''bash
docker compose up -d
'''
Then access the container:
'''bash
docker compose exec taskmaster-dev bash
'''


manager has:
profiles
each profile has an Executor pointer containing all executions
Executor could have all the data about:
- should be restarted?
- Which signal does it stop it?
- Options to discard the program’s stdout/stderr or to redirect them to files
- Environment variables to set before launching the program
- A working directory to set before launching the program
- An umask to set before launching the program
- Whether to start this program at launch or not (?)