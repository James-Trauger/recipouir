package com.jamestrauger.recipouir.models;

import java.util.List;

import jakarta.persistence.CascadeType;
import jakarta.persistence.Column;
import jakarta.persistence.Entity;
import jakarta.persistence.FetchType;
import jakarta.persistence.GeneratedValue;
import jakarta.persistence.GenerationType;
import jakarta.persistence.Id;
import jakarta.persistence.OneToMany;
import jakarta.persistence.Table;

@Entity
@Table(name  = "users") // postgres has default table 'user' already
public class User {
    
    @Id
    @GeneratedValue(strategy = GenerationType.AUTO)
    private Long id;
    @Column(unique = true, length = 31, nullable = false)
    private String username;
    @Column(length = 31)
    private String firstName;
    @Column(length = 31)
    private String lastName;
    // @OneToMany(cascade = CascadeType.REMOVE, fetch = FetchType.LAZY)
    // private List<Recipe> recipes;

    protected User() {}

    public User(String username, String firstName, String lastName) {
        this.username = username;
        this.firstName = firstName;
        this.lastName = lastName;
    }

    public Long getId() {
        return this.id;
    }

    public String getUsername() {
        return this.username;
    }

    public String getFirstName() {
        return this.firstName;
    }

    public String getLastName() {
        return this.lastName;
    }

    // public List<Recipe> getRecipes() {
    //     return this.recipes;
    // }

    public void setId(Long id) {
        this.id = id;
    }

    public void setUsername(String username) {
        this.username = username;
    }

    public void setFirstName(String firstName) {
        this.firstName = firstName;
    }

    public void setLastName(String lastName) {
        this.lastName = lastName;
    }

    // public void setRecipes(List<Recipe> recipes) {
    //     this.recipes = recipes;
    // }


    @Override
    public boolean equals(Object o) {
        if (this == o) 
            return true;
        if (!(o instanceof User)) {
            return false;
        }

        User user = (User) o;

        return this.id.equals(user.getId()) 
            && this.username.equals(user.username)
            && this.firstName.equals(user.getFirstName())
            && this.lastName.equals(user.getLastName());
    }
}
