package com.jamestrauger.recipouir.models;

import org.hibernate.annotations.ColumnDefault;

import jakarta.persistence.Embeddable;
import jakarta.validation.constraints.Min;

@Embeddable
public record Fraction(
    @ColumnDefault("1")
    @Min(1)
    int numerator,
    @Min(1)
    @ColumnDefault("1")
    int denominator
) {}

/* 
import java.io.Serializable;
import java.util.Objects;

import jakarta.persistence.Embeddable;

@Embeddable
public class Fraction implements Serializable {
    private int numerator;
    private int denominator;

    public Fraction(int num, int denom) {
        this.numerator = num;
        this.denominator = denom;
    }

    public int num() {
        return this.numerator;
    }

    public int denom() {
        return this.denominator;
    }

    public void setNum(int num) {
        this.numerator = num;
    }

    public void setDenom(int denom) {
        this.denominator = denom;
    }

    @Override
    public boolean equals(Object o) {
        if (this == o)
            return true;
        if (!(o instanceof Fraction))
            return false;
        
        Fraction frac = (Fraction) o;

        return this.numerator == frac.numerator && this.denominator == frac.denominator;
    }

    @Override
    public int hashCode() {
        //works for object composition
        return Objects.hash(this.numerator, this.denominator);
    }
}
*/