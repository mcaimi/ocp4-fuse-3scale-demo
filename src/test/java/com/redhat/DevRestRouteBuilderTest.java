package com.redhat;

import org.springframework.boot.test.context.SpringBootTest;
import org.junit.Test;
import org.springframework.test.context.ActiveProfiles;
import org.springframework.beans.factory.annotation.Autowired;
import org.junit.runner.RunWith;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.boot.test.context.SpringBootTest.WebEnvironment;
import org.springframework.test.context.junit4.SpringRunner;
import org.springframework.boot.test.web.client.TestRestTemplate;
import org.springframework.http.ResponseEntity;
import org.springframework.http.HttpStatus;
import org.springframework.http.HttpHeaders;
import org.springframework.http.MediaType;
import java.util.concurrent.TimeUnit;
import org.springframework.http.HttpEntity;
import static org.assertj.core.api.Assertions.assertThat;
import java.util.HashMap;
import java.util.Map;

import java.util.Random;
import java.util.UUID;

@ActiveProfiles("test")
@RunWith(SpringRunner.class)
@SpringBootTest(webEnvironment = WebEnvironment.RANDOM_PORT)
public class DevRestRouteBuilderTest {

    @Autowired
    private TestRestTemplate restTemplate;

    @Test
    public void testPublishEndpoint() {
        // publish a new entity
        String randomUserid = UUID.randomUUID().toString();
        String randomTitle = UUID.randomUUID().toString();
        String randomBody = UUID.randomUUID().toString();

        // create a new request header
        HttpHeaders testHeader = new HttpHeaders();
        // set content type to JSON
        testHeader.setContentType(MediaType.APPLICATION_JSON);

        // build request payload
        Map<String, Object> jsonPayload = new HashMap<>();
        jsonPayload.put("userId", randomUserid);
        jsonPayload.put("title", randomTitle);
        jsonPayload.put("body", randomBody);
        HttpEntity<Map<String, Object>> entity = new HttpEntity<>(jsonPayload, testHeader);

        // send POST request
        ResponseEntity<String> response = restTemplate.postForEntity("/camel/api/post", entity, String.class);
        assertThat(response.getStatusCode()).isEqualTo(HttpStatus.NO_CONTENT);
    }

    @Test
    public void testPrometheusEndpoint() {
        ResponseEntity<String> response = restTemplate.getForEntity("/actuator/prometheus", String.class);
        assertThat(response.getStatusCode()).isEqualTo(HttpStatus.OK);
    }
}