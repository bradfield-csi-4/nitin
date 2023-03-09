-- noinspection SqlNoDataSourceInspectionForFile

/* (re)create tables */
create table if not exists department
(
    dep_id serial primary key,
    name   varchar(100)
);
create table if not exists employee
(
    emp_id     serial primary key,
    dep_id     integer references department (dep_id),
    manager_id integer,
    name       varchar(100),
    salary     integer
);
create table if not exists bonus
(
    bonus_id serial primary key,
    emp_id   integer references employee (emp_id),
    amount   integer,
    time     timestamp
);

/* load data into tables */
insert into department (name)
values ('sales'),
       ('marketing'),
       ('engineering');

insert into employee (dep_id, manager_id, name, salary) (select ('{1,1,1,2,3}'::int[])[floor(random() * 5) + 1], /* skew towards sales */
                                                                1 + n % 10, /* first 10 employees manage all the others */
                                                                md5(random()::text), /* name is just a made up string */
                                                                min_salary + random() * (max_salary - min_salary) /* uniform dist of salaries in range */
                                                         from generate_series(1, num_employees) as n);


