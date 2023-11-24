package ru.marlin.example.kafka.consumer;

import lombok.extern.slf4j.Slf4j;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.kafka.annotation.EnableKafka;
import org.springframework.kafka.annotation.KafkaListener;
import org.springframework.kafka.support.Acknowledgment;

@SpringBootApplication
@EnableKafka
@Slf4j
public class MarlinKafkaConsumerApplication {

    public static void main(String[] args) {
        SpringApplication.run(MarlinKafkaConsumerApplication.class, args);
    }


    @KafkaListener(topics = "${topic}")
    public void listenSprutTestStatus(String count, Acknowledgment ack) {
        log.info("Message with count={} received", count);
        ack.acknowledge();
    }
}
