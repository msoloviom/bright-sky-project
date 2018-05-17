package com.ds.canopus.controller;

import com.ds.canopus.domain.Investor;
import com.ds.canopus.resource.PostResource;
import com.ds.canopus.service.InvestorService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpEntity;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import java.util.List;

import static org.springframework.util.CollectionUtils.isEmpty;

@RestController
@RequestMapping("/api")
public class InvestorController {

    private static final String SUCCESSFULLY_CREATED = "Investor was successfully created !";

    @Autowired
    private InvestorService investorService;

    @GetMapping("/investors")
    public HttpEntity<List<Investor>> getAllInvestors() {
        List<Investor> investors = investorService.getAllInvestors();
        if(isEmpty(investors)) {
            return new ResponseEntity<>(HttpStatus.NOT_FOUND);
        } else {
            return new ResponseEntity<>(investors, HttpStatus.OK);
        }
    }

    @GetMapping("/investors/{id}")
    public HttpEntity<Investor> getInvestorById(@PathVariable Long id) {
        Investor investor = investorService.getInvestorById(id);
        if (investor == null) {
            return new ResponseEntity<>(HttpStatus.NOT_FOUND);
        } else {
            return new ResponseEntity<>(investor, HttpStatus.OK);
        }
    }

    @PostMapping("/investors")
    public HttpEntity<PostResource> saveInvestor(@RequestBody Investor investor) {
        Long investorId = investorService.saveInvestor(investor);
        return new ResponseEntity<>(new PostResource(investorId, SUCCESSFULLY_CREATED), HttpStatus.CREATED);
    }

}