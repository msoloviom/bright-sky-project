package com.ds.canopus.mq;

import com.ds.canopus.domain.Investor;
import com.ds.canopus.service.InvestorService;
import org.springframework.amqp.rabbit.annotation.RabbitListener;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

@Service
public class SaverMqListener {

    @Autowired
    private InvestorService investorService;

    @RabbitListener(queues = "investor")
    public void consumeInvestorForSaving(final Investor investor) {
        Long investorId = investorService.saveInvestor(investor);
        System.out.println("New investor was saved. Generated Id: " + investorId);
    }

}
