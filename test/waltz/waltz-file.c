
#include "o5c.h"
#include <math.h>

#define PI		3.1415927
#define MOD_NUM		100
#define get_y(val)	(val % MOD_NUM)
#define get_x(val)	((int) (val / MOD_NUM))


int p1_offset, p2_offset, p3_offset;
StrID tee_str_id, fork_str_id, arrow_str_id;

BOOLEAN first = TRUE;

void
init()
{
	StringEntry str_entry;

	p1_offset = dollar_litbind("p1");
	p2_offset = dollar_litbind("p2");
	p3_offset = dollar_litbind("p3");

	str_entry = AddString("tee",attr_values);
	tee_str_id = str_entry->str_id;
	str_entry = AddString("fork",attr_values);
	fork_str_id = str_entry->str_id;
	str_entry = AddString("arrow",attr_values);
	arrow_str_id = str_entry->str_id;

	first = FALSE;
}


/**********************************************************************
 * This function is passed two points and calculates the angle between the
 * line defined by these points and the x-axis.
 **********************************************************************/
float
get_angle(p1,p2)
int p1,p2;
{
	int delta_x, delta_y;

	/* Calculate (x2 - x1) and (y2 - y1).  The points are passed in the
	 * form x1y1 and x2y2.  get_x() and get_y() are passed these points
	 * and return the x and y values respectively.  For example,
	 * get_x(1020) returns 10. */
	delta_x = get_x(p2) - get_x(p1);
	delta_y = get_y(p2) - get_y(p1);

	if (delta_x == 0) {
		if (delta_y > 0)
			return(PI/2);
		else if (delta_y < 0)
			return(-PI/2);
	}
	else if (delta_y == 0) {
		if (delta_x > 0)
			return(0.0);
		else if (delta_x < 0)
			return(PI);
	}
	else
		return((float) atan2((double) delta_y,(double) delta_x));
}


/**********************************************************************
 * This procedure is passed the basepoint of the intersection of two lines
 * as well as the other two endpoints of the lines and calculates the
 * angle inscribed by these three points.
 **********************************************************************/
float
inscribed_angle(basepoint,p1,p2)
int basepoint, p1, p2;
{
	float angle1, angle2, temp;

	/* Get the angle between line #1 and the origin and the angle
	 * between line #2 and the origin, and then subtract these values. */
	angle1 = get_angle(basepoint,p1);
	angle2 = get_angle(basepoint,p2);
	temp = angle1 - angle2;
	if (temp < 0.0)
		temp = -temp;

	/* We always want the smaller of the two angles inscribed, so if the
	 * answer is greater than 180 degrees, calculate the smaller angle and
	 * return it. */
	if (temp > PI)
		temp = 2*PI - temp;
	if (temp < 0.0)
		return(-temp);
	return(temp);
}


void
make_3_junction(argc)
int argc;
{
	int basepoint,p1,p2,p3;
	int shaft,barb1,barb2;
	float angle12, angle13, angle23;
	float sum, sum1213, sum1223, sum1323;
	float delta;

	if (first)
		init();

	basepoint = args[0].data.ival;
	p1 = args[1].data.ival;
	p2 = args[2].data.ival;
	p3 = args[3].data.ival;

	angle12 = inscribed_angle(basepoint,p1,p2);
	angle13 = inscribed_angle(basepoint,p1,p3);
	angle23 = inscribed_angle(basepoint,p2,p3);

	sum1213 = angle12 + angle13;
	sum1223 = angle12 + angle23;
	sum1323 = angle13 + angle23;

	if (sum1213 < sum1223) {
		if (sum1213 < sum1323) {
			sum = sum1213;
			shaft = p1; barb1 = p2; barb2 = p3;
		}
		else {
			sum = sum1323;
			shaft = p3; barb1 = p1; barb2 = p2;
		}
	}
	else {
		if (sum1223 < sum1323) {
			sum = sum1223;
			shaft = p2; barb1 = p1; barb2 = p3;
		}
		else {
			sum = sum1323;
			shaft = p3; barb1 = p1; barb2 = p2;
		}
	}

	delta = sum - PI;
	if (delta < 0.0)
		delta = -delta;

	if (delta < 0.001)
		dollar_str_val(tee_str_id);
	else if (sum > PI)
		dollar_str_val(fork_str_id);
	else
		dollar_str_val(arrow_str_id);

	dollar_tab1(p1_offset);
	dollar_int_val(barb1);

	dollar_tab1(p2_offset);
	dollar_int_val(shaft);

	dollar_tab1(p3_offset);
	dollar_int_val(barb2);
}
