all: 
					docker compose up -d --build
					docker compose logs -f bot

stop:
					docker compose stop

down:
					docker compose down

clean:
					docker image prune -f

fclean:				down
					docker image prune -f -a

re:					fclean all

image:
					docker image ls

.PHONY: all stop down clean fclean re image
