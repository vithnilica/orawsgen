package orawsgen;

import java.io.IOException;
import java.io.InputStream;
import java.io.InputStreamReader;
import java.io.Reader;
import java.io.StringWriter;
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
import javax.xml.transform.stream.StreamResult;
import javax.xml.transform.stream.StreamSource;

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
		if(x==null || x.xmlstr==null)return null;
		try {
			SQLXML sqlXml = con.createSQLXML();
			
			StreamResult strRes=sqlXml.setResult(StreamResult.class);		
			
			//strRes.getOutputStream().write(x.xml.getBytes("UTF8"));
			strRes.getOutputStream().write(x.xmlstr.getBytes());
			return sqlXml;
		} catch (IOException e) {
			throw new SQLException(e);
		}		

	}
	
	public static XmlAny convertSqlXml2XmlAny(SQLXML sqlXml)throws SQLException {
		XmlAny x=new XmlAny();
		if(sqlXml==null)return x;
		try {
			StreamSource strSrc=sqlXml.getSource(StreamSource.class);
			InputStream is=strSrc.getInputStream();
			
			final int bufferSize = 1024;
			final char[] buffer = new char[bufferSize];
			final StringBuilder out = new StringBuilder();
			//Reader in = new InputStreamReader(inputStream, "UTF-8");
			Reader in = new InputStreamReader(is);
			for (; ; ) {
			    int rsz = in.read(buffer, 0, buffer.length);
			    if (rsz < 0)
			        break;
			    out.append(buffer, 0, rsz);
			}
			x.xmlstr=out.toString();

			return x;
		} catch (IOException e) {
			throw new SQLException(e);
		}	

	}




}
