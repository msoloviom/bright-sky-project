package com.ds.canopus.repository;

import com.ds.canopus.domain.Investor;
import org.springframework.data.repository.CrudRepository;
import org.springframework.stereotype.Repository;

import java.util.List;

@Repository
public interface InvestorRepository extends CrudRepository <Investor, Long> {

    List<Investor> findAll();

    Investor findInvestorById(Long id);

    Investor save(Investor investor);
}
