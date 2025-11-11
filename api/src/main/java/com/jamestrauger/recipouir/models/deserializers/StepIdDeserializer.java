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

import com.jamestrauger.recipouir.models.StepId;
import com.jamestrauger.recipouir.models.Recipe;

public class StepIdDeserializer extends StdDeserializer<StepId> {


    public StepIdDeserializer() {
        this(null);
    }

    public StepIdDeserializer(Class<?> vc) {
        super(vc);
    }

    @Override
    public StepId deserialize(JsonParser jp, DeserializationContext ctxt)
            throws IOException, JacksonException {
        
        JsonNode node = jp.getCodec().readTree(jp);
        int number = node.get("number").asInt();
        
        Recipe recipe = null;
    
        // parent Step
        JsonStreamContext StepContext = jp.getParsingContext().getParent();
        if (StepContext != null) {
            // parent recipe
            JsonStreamContext recipeContext = StepContext.getParent();
            if (recipeContext != null) {
                // set the recipe reference to the parent
                recipe = (Recipe) recipeContext.getCurrentValue();
            }
        }

        return new StepId(recipe, number);
    }
}