package com.ds.saverMQ.producer;


import java.io.Serializable;

public class Investor implements Serializable {

    private Long id;

    private String name;

    private String cert;

    public Long getId() {
        return id;
    }

    public void setId(Long id) {
        this.id = id;
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public String getCert() {
        return cert;
    }

    public void setCert(String cert) {
        this.cert = cert;
    }
}
