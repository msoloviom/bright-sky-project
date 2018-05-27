package com.ds.saverMQ.producer;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpEntity;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.Mapping;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import java.io.IOException;
import java.util.concurrent.TimeoutException;

@RestController
@RequestMapping("/api")
public class SaverController {

    @Autowired
    private ProducerMQ producerMQ;

    @PostMapping("/mq/investor")
    public HttpEntity<?> putInvestorSavingInQueue(@RequestBody Investor investor) throws IOException, TimeoutException {
        producerMQ.produce(investor);
        return new ResponseEntity<>(HttpStatus.ACCEPTED);
    }
}
