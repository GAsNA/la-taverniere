all: 
					docker compose up --build

clean:
					docker compose stop
					docker compose down

fclean: 			clean
					docker system prune -f --all

re:					fclean all

image:
					docker image ls

.PHONY: all clean fclean re image
