package com.jamestrauger.recipouir.models.deserializers;

import java.io.IOException;

import com.fasterxml.jackson.core.JacksonException;
import com.fasterxml.jackson.databind.DeserializationContext;
import com.fasterxml.jackson.databind.JsonDeserializer;
import com.fasterxml.jackson.databind.JsonMappingException;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.deser.ResolvableDeserializer;
import com.fasterxml.jackson.databind.deser.std.StdDeserializer;
import com.fasterxml.jackson.databind.node.ObjectNode;
import com.fasterxml.jackson.core.JsonParser;
import com.fasterxml.jackson.core.JsonStreamContext;

import com.jamestrauger.recipouir.models.IngredientId;
import com.jamestrauger.recipouir.models.Recipe;

public class IngredientIdDeserializer extends StdDeserializer<IngredientId> {


    public IngredientIdDeserializer() {
        this(null);
    }

    public IngredientIdDeserializer(Class<?> vc) {
        super(vc);
    }

    @Override
    public IngredientId deserialize(JsonParser jp, DeserializationContext ctxt)
            throws IOException, JacksonException {
        
        JsonNode node = jp.getCodec().readTree(jp);
        String name = node.get("name").asText();
        
        Recipe recipe = null;
    
        // parent ingredient
        JsonStreamContext ingredientContext = jp.getParsingContext().getParent();
        if (ingredientContext != null) {
            // parent recipe
            JsonStreamContext recipeContext = ingredientContext.getParent();
            if (recipeContext != null) {
                // set the recipe reference to the parent
                recipe = (Recipe) recipeContext.getCurrentValue();
            }
        }

        return new IngredientId(name, recipe);
    }
}