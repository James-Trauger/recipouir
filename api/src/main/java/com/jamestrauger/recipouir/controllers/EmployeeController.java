package com.jamestrauger.recipouir.controllers;

import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RestController;

import com.jamestrauger.recipouir.models.EmployeeRepository;
import com.jamestrauger.recipouir.models.Employee;
import com.jamestrauger.recipouir.errors.NotFoundException;

@RestController
public class EmployeeController {
    
    private final  EmployeeRepository employeeRepository;

    public EmployeeController(EmployeeRepository employeeRepository) {
        this.employeeRepository = employeeRepository;
    }

    @GetMapping("/employees")
    public Iterable<Employee> findAllEmployees() {
        return this.employeeRepository.findAll();
    }

    @PostMapping("/employees")
    public Employee newEmployee(@RequestBody Employee employee) {
        return this.employeeRepository.save(employee);
    }

    @GetMapping("/employees/{id}")
    public Employee one(@PathVariable Long id) {
        return this.employeeRepository.findById(id)
            .orElseThrow(() -> new NotFoundException(id));
    }
}
