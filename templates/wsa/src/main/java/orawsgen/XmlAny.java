package orawsgen;

import javax.xml.bind.annotation.XmlAccessType;
import javax.xml.bind.annotation.XmlAccessorType;
import javax.xml.bind.annotation.XmlAnyElement;
import javax.xml.bind.annotation.XmlType;

import org.w3c.dom.Element;

//jedno xml v tagu any
@XmlAccessorType(XmlAccessType.FIELD)
@XmlType(name = "", propOrder = {
    "any"
})
public class XmlAny {
    @XmlAnyElement(lax = true)
    protected Element any;

    public Element getAny() {
        return any;
    }

    public void setAny(Element value) {
        this.any = value;
    }
}
