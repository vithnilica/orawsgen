package orawsgen;

import java.sql.Array;
import java.sql.Connection;
import java.sql.SQLException;
import java.sql.SQLXML;
import java.sql.Struct;
import java.sql.Timestamp;
import java.util.Arrays;
import java.util.Date;
import java.util.GregorianCalendar;

import javax.xml.transform.dom.DOMResult;
import javax.xml.transform.dom.DOMSource;

import org.w3c.dom.Element;

import oracle.jdbc.OracleConnection;
import oracle.sql.StructDescriptor;

public class DBUtils {

	public static Struct createStruct(Connection con, String typeName, Object[] arr) throws SQLException {
		//ohejbak ktery zajisti fungovani i po pridani nove polozky do typu
		StructDescriptor sd=StructDescriptor.createDescriptor(typeName, con);
		int newLength=sd.getLength();
		if(arr.length>newLength) {
			//v databazi ma typ mene polozek nez s cim pocita java
			throw new SQLException("v databazi ma typ mene polozek ("+typeName+", "+arr.length+", "+newLength);
		}else if(arr.length<newLength) {
			//TODO jen zalogovat
		}
		arr=Arrays.copyOf(arr, newLength);
		return con.createStruct(typeName, arr);
	}

	public static Array createArray(Connection con, String typeName, Object[] arr)throws SQLException {
		Array sqlArr = ((OracleConnection)con).createOracleArray(typeName,arr);
		return sqlArr;
	}

	/*
	public static XMLGregorianCalendar convertSqlTs2XmlCal(Timestamp ts)throws SQLException {
		if(ts==null)return null;
		GregorianCalendar cal1 = new GregorianCalendar();
		cal1.setTimeInMillis(ts.getTime());
		try {
			return DatatypeFactory.newInstance().newXMLGregorianCalendar(cal1);
		} catch (DatatypeConfigurationException e) {
			throw new SQLException(e);
		}
	}

	public static Timestamp convertXmlCal2SqlTs(XMLGregorianCalendar cal)throws SQLException {
		if(cal==null)return null;
		return new Timestamp(cal.toGregorianCalendar().getTimeInMillis());
	}
	*/

	public static Date convertSqlTs2Date(Timestamp ts)throws SQLException {
		if(ts==null)return null;
		GregorianCalendar cal1 = new GregorianCalendar();
		cal1.setTimeInMillis(ts.getTime());
		return cal1.getTime();
	}

	public static Timestamp convertDate2SqlTs(Date d)throws SQLException {
		if(d==null)return null;
		return new Timestamp(d.getTime());

	}

	
	public static SQLXML convertXmlAny2SqlXml(Connection con, XmlAny x)throws SQLException {
		if(x==null)return null;
		SQLXML sqlXml = con.createSQLXML();
		Element el=x.getAny();
		if(el!=null) {
			DOMResult domRes=sqlXml.setResult(DOMResult.class);
			domRes.setNode(el);
			return sqlXml;
		}else {
			return null;
		}
	}
	
	public static XmlAny convertSqlXml2XmlAny(SQLXML sqlXml)throws SQLException {
		XmlAny x=new XmlAny();
		if(sqlXml==null)return x;
		DOMSource domSrc=sqlXml.getSource(DOMSource.class);
		x.setAny((Element) domSrc.getNode());
		return x;
	}




}
