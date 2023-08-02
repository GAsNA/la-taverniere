all: 
					docker compose up -d --build
					docker compose logs -f bot

stop:
					docker compose stop

clean:
					docker compose down

fclean: 			clean
					docker image prune -f

re:					fclean all

image:
					docker image ls

.PHONY: all stop clean fclean re image
