# frontend/deployments/terraform/main.tf
provider "aws" {
  region = "us-east-1"
}

resource "aws_ecr_repository" "backend" {
  name = "autosysadmin/backend"
}

resource "aws_ecr_repository" "frontend" {
  name = "autosysadmin/frontend"
}

resource "aws_ecs_cluster" "autosysadmin" {
  name = "autosysadmin-cluster"
}

resource "aws_ecs_task_definition" "backend" {
  family                   = "autosysadmin-backend"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = 512
  memory                   = 1024
  execution_role_arn       = aws_iam_role.ecs_task_execution_role.arn

  container_definitions = jsonencode([{
    name      = "backend"
    image     = "${aws_ecr_repository.backend.repository_url}:latest"
    cpu       = 512
    memory    = 1024
    essential = true
    portMappings = [{
      containerPort = 8080
      hostPort      = 8080
    }]
    environment = [
      { name = "DB_HOST", value = aws_db_instance.autosysadmin.address },
      { name = "DB_USER", value = "autosysadmin" },
      { name = "DB_PASSWORD", value = var.db_password },
      { name = "DB_NAME", value = "autosysadmin" },
      { name = "JWT_SECRET", value = var.jwt_secret }
    ]
  }])
}

resource "aws_ecs_service" "backend" {
  name            = "autosysadmin-backend"
  cluster         = aws_ecs_cluster.autosysadmin.id
  task_definition = aws_ecs_task_definition.backend.arn
  desired_count   = 2
  launch_type     = "FARGATE"

  network_configuration {
    subnets          = aws_subnet.private.*.id
    security_groups  = [aws_security_group.ecs_tasks.id]
    assign_public_ip = false
  }

  load_balancer {
    target_group_arn = aws_lb_target_group.backend.arn
    container_name   = "backend"
    container_port   = 8080
  }
}

resource "aws_lb" "autosysadmin" {
  name               = "autosysadmin-lb"
  internal           = false
  load_balancer_type = "application"
  subnets            = aws_subnet.public.*.id
  security_groups    = [aws_security_group.lb.id]
}

resource "aws_lb_target_group" "backend" {
  name        = "autosysadmin-backend"
  port        = 80
  protocol    = "HTTP"
  target_type = "ip"
  vpc_id      = aws_vpc.main.id

  health_check {
    path = "/health"
  }
}

resource "aws_lb_listener" "backend" {
  load_balancer_arn = aws_lb.autosysadmin.arn
  port              = "80"
  protocol          = "HTTP"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.backend.arn
  }
}