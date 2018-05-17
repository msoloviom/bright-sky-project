package com.ds.canopus.service;

import com.ds.canopus.domain.Investor;
import com.ds.canopus.repository.InvestorRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.List;

@Service
public class InvestorServiceImpl implements InvestorService {

    @Autowired
    private InvestorRepository investorRepository;

    @Override
    public List<Investor> getAllInvestors() {
        return investorRepository.findAll();
    }

    @Override
    public Investor getInvestorById(Long id) {
        return investorRepository.findInvestorById(id);
    }

    @Override
    public Long saveInvestor(Investor investor) {
        return investorRepository.save(investor).getId();
    }
}
