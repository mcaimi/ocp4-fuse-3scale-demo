package com.redhat;

import java.sql.Timestamp;
import java.time.LocalDateTime;
import java.util.Map;

import org.apache.camel.Exchange;
import org.apache.camel.Processor;
import org.apache.camel.builder.RouteBuilder;
import org.apache.camel.model.dataformat.JsonLibrary;
import org.apache.camel.model.rest.RestBindingMode;
import org.apache.camel.util.toolbox.AggregationStrategies;
import org.apache.camel.util.toolbox.FlexibleAggregationStrategy;
import org.springframework.stereotype.Component;
import org.springframework.context.annotation.Profile;

import com.redhat.db.PostEntity;

@Component
@Profile("test")
public class DevRestRouteBuilder extends RouteBuilder {

    @Override
    public void configure() {

        onException(Exception.class)
                .handled(true)
                .setHeader(Exchange.HTTP_RESPONSE_CODE, constant(204))
                .setBody().constant("Invalid request")
                .marshal().json(JsonLibrary.Jackson)
                .stop();

        from("timer://test?repeatCount=1")
                .setBody(simple("CREATE TABLE post_entity (user_id INTEGER NOT NULL, body VARCHAR(255), title VARCHAR(255), created TIMESTAMP DEFAULT CURRENT_TIMESTAMP());"))
                .log("------------------------------")
                .log("CREATE : ${body}")
                .log("------------------------------")
                .to("jdbc:datasource")
                .log("Table created.");

        restConfiguration()
                .enableCORS(true)
                .contextPath("/camel")
                .apiContextPath("/api-doc")
                .apiProperty("api.title", "Fuse REST to JDBC PoC")
                .apiProperty("api.version", "v1")
                .apiContextRouteId("doc-api")
                .component("servlet")
                .bindingMode(RestBindingMode.json);

        rest("/api/post")
			.description("Posts service")
	        .post()
	            .description("Inserts a new post")
	            .type(PostEntity.class)
	            .route()
	            	.routeId("api-process")
	                .to("direct:insert")
	            .endRest()
			.get()
				.description("Lists available posts")
				.outType(PostEntity[].class)
				.route()
					.routeId("post-select")
					.to("direct:select")
				.endRest();

        from("direct:insert")
                .log("------------------------------")
                .log("INSERT : ${body}")
                .log("------------------------------")
                .setBody(simple("INSERT INTO post_entity (user_id,body,title) VALUES (${body.userId},'${body.body}','${body.title}')"))
                .to("jdbc:datasource")
                .log("Entity successfully saved.");

        from("direct:select")
                .log("------------------------------")
                .log("SELECT")
                .log("------------------------------")
                .setBody(constant("select * from post_entity"))
                .to("jdbc:dataSource")
                .split(body()).aggregationStrategy(AggregationStrategies.groupedBody())
                    .process(exchange -> {
                        @SuppressWarnings("unchecked")
                        Map<String, Object> row = (Map<String, Object>) exchange.getIn().getBody();
                        Integer userId = (Integer) row.get("USER_ID");
                        String title = (String) row.get("TITLE");
                        String body = (String) row.get("BODY");
                        LocalDateTime created = ((Timestamp) row.get("CREATED")).toLocalDateTime();
                        PostEntity entity = new PostEntity(userId, title, body, created);
                        exchange.getMessage().setBody(entity);
                    })
                .end()
                .log("Result as DTO : ${body}");
    }
}
