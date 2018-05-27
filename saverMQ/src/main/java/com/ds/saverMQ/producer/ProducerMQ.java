package com.ds.saverMQ.producer;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.amqp.core.AmqpTemplate;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;

import java.io.IOException;
import java.util.concurrent.TimeoutException;

@Component
public class ProducerMQ {

    public final static String QUEUE_NAME = "investor_q";

    private static final Logger log = LoggerFactory.getLogger(ProducerMQ.class);

    @Autowired
    private AmqpTemplate amqpTemplate;

    @Value("${jsa.rabbitmq.exchange}")
    private String exchange;

    public void produce(Investor investor){
        amqpTemplate.convertAndSend(exchange, "", investor);
        System.out.println("Send msg = " + investor);
    }

    //@RabbitListener(queues = QUEUE_NAME)
    public void publishInvestor(Investor investor) throws IOException, TimeoutException {
        //rabbitTemplate.setQueue(QUEUE_NAME);
        log.info("recording: " + investor);
        //rabbitTemplate.convertAndSend("", "", "ljgfljgf");
        /*Connection connection = createRemoteConnection();
        Channel channel = connection.createChannel();

        channel.queueDeclare(QUEUE_NAME, false, false, false, null);
        channel.basicPublish("", QUEUE_NAME, null, investor.);
        System.out.println(" [x] Sent '" + message + "'");

        channel.close();
        connection.close();*/
    }
}
