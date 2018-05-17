package com.ds.canopus.service;

import com.ds.canopus.domain.Investor;

import java.util.List;

public interface InvestorService {
    List<Investor> getAllInvestors();
    Investor getInvestorById(Long id);
    Long saveInvestor(Investor investor);
}
