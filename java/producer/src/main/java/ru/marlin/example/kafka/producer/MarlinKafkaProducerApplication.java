package ru.marlin.example.kafka.producer;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.kafka.annotation.EnableKafka;
import org.springframework.kafka.core.KafkaTemplate;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RestController;

import java.util.concurrent.atomic.AtomicLong;

@SpringBootApplication
@EnableKafka
public class MarlinKafkaProducerApplication {

    public static void main(String[] args) {
        SpringApplication.run(MarlinKafkaProducerApplication.class, args);
    }


    @RestController
    static public class KafkaController {

        private final AtomicLong counter = new AtomicLong();
        private final KafkaTemplate<String, Object> template;
        private final String topic;

        public KafkaController(KafkaTemplate<String, Object> template,
                               @Value("${topic}") String topic) {
            this.template = template;
            this.topic = topic;
        }

        @PostMapping(path="kafka")
        public void sendMessage() {
            var count = String.valueOf(counter.getAndIncrement());
            template.send(topic, count, count);
        }
    }
}
